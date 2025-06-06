#!/bin/bash
# 混沌工程故障诊断脚本
# 用于诊断 AZ 网络故障和 Pod 健康状况

set -euo pipefail

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
NAMESPACE=${NAMESPACE:-"default"}
RELEASE_NAME=${RELEASE_NAME:-""}
CHECK_NETWORK=${CHECK_NETWORK:-"true"}

echo_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_header() {
    echo "=============================================="
    echo "🔍 混沌工程故障诊断工具"
    echo "时间: $(date)"
    echo "命名空间: $NAMESPACE"
    echo "=============================================="
}

check_cronjob_status() {
    echo_info "检查 CronJob 状态..."
    
    if [ -n "$RELEASE_NAME" ]; then
        CRONJOB_NAME="${RELEASE_NAME}-rebalance"
    else
        CRONJOB_NAME=$(kubectl get cronjobs -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}' | grep -E "rebalance|chaos" | head -1)
    fi
    
    if [ -z "$CRONJOB_NAME" ]; then
        echo_error "未找到重新平衡 CronJob"
        return 1
    fi
    
    echo "CronJob 名称: $CRONJOB_NAME"
    kubectl get cronjob $CRONJOB_NAME -n $NAMESPACE
    
    echo -e "\n最近的 Job 执行:"
    kubectl get jobs -n $NAMESPACE --sort-by=.metadata.creationTimestamp | grep $CRONJOB_NAME | tail -3
    
    echo -e "\n最近 Job 的日志:"
    LATEST_JOB=$(kubectl get jobs -n $NAMESPACE --sort-by=.metadata.creationTimestamp -o jsonpath='{.items[-1:].metadata.name}' | grep $CRONJOB_NAME || echo "")
    if [ -n "$LATEST_JOB" ]; then
        kubectl logs -l job-name=$LATEST_JOB -n $NAMESPACE --tail=50
    fi
}

check_pod_distribution() {
    echo_info "检查 Pod 在各 AZ 的分布..."
    
    echo -e "\n按 AZ 分组的 Pod 分布:"
    kubectl get pods -n $NAMESPACE -o custom-columns="NAME:.metadata.name,NODE:.spec.nodeName,AZ:.metadata.labels.topology\.kubernetes\.io/zone,STATUS:.status.phase,READY:.status.conditions[?(@.type=='Ready')].status" | \
        awk 'NR>1 {az[$3]++; if($5=="True") ready[$3]++} END {for(i in az) printf "AZ: %-20s Total: %d, Ready: %d, Ratio: %.2f%%\n", i, az[i], ready[i]+0, (ready[i]+0)/az[i]*100}'
    
    echo -e "\n详细 Pod 状态:"
    kubectl get pods -n $NAMESPACE -o wide --sort-by=.spec.nodeName
}

check_node_status() {
    echo_info "检查节点状态..."
    
    echo "节点按 AZ 分组:"
    kubectl get nodes -o custom-columns="NAME:.metadata.name,AZ:.metadata.labels.topology\.kubernetes\.io/zone,STATUS:.status.conditions[?(@.type=='Ready')].status,TAINTS:.spec.taints[*].key" | \
        awk 'NR>1 {az[$2]++; if($3=="True") ready[$2]++} END {for(i in az) printf "AZ: %-20s Total: %d, Ready: %d\n", i, az[i], ready[i]+0}'
    
    echo -e "\n节点详细信息:"
    kubectl get nodes -o wide
}

check_unhealthy_pods() {
    echo_info "检查不健康的 Pod..."
    
    local unhealthy_count=0
    
    echo "检查未就绪的 Pod:"
    kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running,status.phase!=Succeeded
    
    echo -e "\n检查重启次数过多的 Pod (>5次):"
    kubectl get pods -n $NAMESPACE -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.spec.nodeName}{" "}{.status.containerStatuses[0].restartCount}{" "}{.metadata.labels.topology\.kubernetes\.io/zone}{"\n"}{end}' | \
        awk '$3 > 5 {print "Pod:", $1, "Node:", $2, "Restarts:", $3, "AZ:", $4; unhealthy++} END {if(unhealthy>0) printf "\n发现 %d 个重启过多的 Pod\n", unhealthy}'
    
    echo -e "\n检查长时间未就绪的 Pod (>5分钟):"
    kubectl get pods -n $NAMESPACE -o json | jq -r '.items[] | select(.status.conditions[]? | select(.type=="Ready" and .status=="False")) | select((now - (.metadata.creationTimestamp | fromdateiso8601)) > 300) | "\(.metadata.name) \(.spec.nodeName) \((now - (.metadata.creationTimestamp | fromdateiso8601))/60 | floor)min \(.metadata.labels."topology.kubernetes.io/zone" // "unknown")"' | \
        while read line; do
            if [ -n "$line" ]; then
                echo "Pod: $line"
                ((unhealthy_count++))
            fi
        done
    
    if [ $unhealthy_count -eq 0 ]; then
        echo_success "没有发现明显不健康的 Pod"
    else
        echo_warn "发现 $unhealthy_count 个不健康的 Pod"
    fi
}

check_network_connectivity() {
    if [ "$CHECK_NETWORK" != "true" ]; then
        echo_info "跳过网络连接检查"
        return 0
    fi
    
    echo_info "检查 Pod 网络连接..."
    
    local failed_pods=0
    
    # 获取运行中的 Pod
    kubectl get pods -n $NAMESPACE -o jsonpath='{range .items[?(@.status.phase=="Running")]}{.metadata.name}{" "}{.spec.nodeName}{" "}{.metadata.labels.topology\.kubernetes\.io/zone}{"\n"}{end}' | \
        while read pod_name node_name az; do
            if [ -n "$pod_name" ]; then
                echo -n "检查 $pod_name (AZ: $az): "
                if kubectl exec $pod_name -n $NAMESPACE -- timeout 5 wget -q --spider http://kubernetes.default.svc.cluster.local 2>/dev/null; then
                    echo_success "网络正常"
                else
                    echo_error "网络异常"
                    ((failed_pods++))
                fi
            fi
        done
    
    if [ $failed_pods -eq 0 ]; then
        echo_success "所有 Pod 网络连接正常"
    else
        echo_warn "发现 $failed_pods 个 Pod 网络连接异常"
    fi
}

check_chaosblade_experiments() {
    echo_info "检查 ChaosBlade 实验状态..."
    
    if command -v blade >/dev/null 2>&1; then
        echo "当前活跃的混沌实验:"
        blade status --type create
    else
        echo_warn "ChaosBlade CLI 未安装，无法查看实验状态"
        
        # 尝试通过 kubectl 查看
        echo "尝试查看 ChaosBlade 相关资源:"
        kubectl get chaosblade -A 2>/dev/null || echo "未找到 ChaosBlade CRD"
    fi
}

check_rbac_permissions() {
    echo_info "检查 RBAC 权限..."
    
    local service_account="default"
    if [ -n "$RELEASE_NAME" ]; then
        service_account="${RELEASE_NAME}"
    fi
    
    echo "检查 ServiceAccount: $service_account"
    kubectl auth can-i delete pods --as=system:serviceaccount:$NAMESPACE:$service_account -n $NAMESPACE
    kubectl auth can-i get nodes --as=system:serviceaccount:$NAMESPACE:$service_account
    kubectl auth can-i exec pods --as=system:serviceaccount:$NAMESPACE:$service_account -n $NAMESPACE
}

generate_recommendations() {
    echo_info "生成建议..."
    
    # 检查 AZ 分布
    local az_count=$(kubectl get nodes -o jsonpath='{.items[*].metadata.labels.topology\.kubernetes\.io/zone}' | tr ' ' '\n' | sort | uniq | wc -l)
    local pod_count=$(kubectl get pods -n $NAMESPACE --field-selector=status.phase=Running | wc -l)
    
    echo "环境概况:"
    echo "- 可用区数量: $az_count"
    echo "- 运行中 Pod 数量: $pod_count"
    
    echo -e "\n建议:"
    
    if [ $az_count -lt 2 ]; then
        echo_warn "建议使用多个可用区以提高高可用性"
    fi
    
    if [ $pod_count -lt 6 ]; then
        echo_warn "Pod 数量较少，建议增加副本数量以验证重新平衡效果"
    fi
    
    echo "- 建议在业务低峰期进行混沌工程测试"
    echo "- 确保有足够的健康节点承载重新调度的工作负载"
    echo "- 监控应用的关键指标和用户体验"
    echo "- 逐步增加故障的影响范围和持续时间"
}

simulate_chaos_test() {
    echo_info "模拟混沌测试建议..."
    
    echo "建议的测试场景:"
    echo "1. 单节点网络故障:"
    echo "   blade create k8s node-network-loss --names <node-name>"
    
    echo -e "\n2. AZ 级别网络故障:"
    local nodes_in_az_a=$(kubectl get nodes -o jsonpath='{.items[?(@.metadata.labels.topology\.kubernetes\.io/zone=="cn-guangzhou-a")].metadata.name}' | tr ' ' ',')
    if [ -n "$nodes_in_az_a" ]; then
        echo "   blade create k8s node-network-loss --names $nodes_in_az_a"
    else
        echo "   blade create k8s node-network-loss --names <node1>,<node2>"
    fi
    
    echo -e "\n3. 手动触发重新平衡测试:"
    if [ -n "$CRONJOB_NAME" ]; then
        echo "   kubectl create job manual-chaos-test --from=cronjob/$CRONJOB_NAME -n $NAMESPACE"
    else
        echo "   kubectl create job manual-chaos-test --from=cronjob/<cronjob-name> -n $NAMESPACE"
    fi
    
    echo -e "\n4. 监控命令:"
    echo "   kubectl get pods -o wide --watch"
    echo "   kubectl logs -f -l job-name=manual-chaos-test"
}

# 主函数
main() {
    print_header
    
    echo_info "开始诊断..."
    echo
    
    check_cronjob_status
    echo
    
    check_node_status
    echo
    
    check_pod_distribution
    echo
    
    check_unhealthy_pods
    echo
    
    check_network_connectivity
    echo
    
    check_chaosblade_experiments
    echo
    
    check_rbac_permissions
    echo
    
    generate_recommendations
    echo
    
    simulate_chaos_test
    echo
    
    echo_success "诊断完成！"
}

# 参数处理
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -r|--release)
            RELEASE_NAME="$2"
            shift 2
            ;;
        --no-network-check)
            CHECK_NETWORK="false"
            shift
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo "选项:"
            echo "  -n, --namespace     指定命名空间 (默认: default)"
            echo "  -r, --release       指定 Helm release 名称"
            echo "  --no-network-check  跳过网络连接检查"
            echo "  -h, --help          显示帮助信息"
            exit 0
            ;;
        *)
            echo_error "未知参数: $1"
            exit 1
            ;;
    esac
done

# 检查 kubectl 命令
if ! command -v kubectl >/dev/null 2>&1; then
    echo_error "kubectl 命令未找到，请先安装 kubectl"
    exit 1
fi

# 检查 jq 命令
if ! command -v jq >/dev/null 2>&1; then
    echo_warn "jq 命令未找到，某些功能可能受限"
fi

main 