import { useRouter } from 'next/router';
import Link from 'next/link';
import styles from '../../../styles/orderDetail.module.css';
import { ordersDemoData } from '../../../utils/demoData';
import Header from '../../../components/Header/Header';

const OrderDetailPage = () => {
  const router = useRouter();
  const { orderId } = router.query;
  const useDemoData = true;

  if (useDemoData && ordersDemoData.data.length > 0) {
    const order = ordersDemoData.data.find(o => o.id === parseInt(orderId as string)) || ordersDemoData.data[0];

    return (
      <>
        <Header />
        <div className={styles.container}>
          <h1 className={styles.title}>订单详情: #{orderId}</h1>
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>订单信息</h2>
            <div className={styles.detailRow}>
              <span className={styles.label}>订单 ID:</span>
              <span className={styles.value}>{order.order_id}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>用户 ID:</span>
              <span className={styles.value}>{order.user_id}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>创建时间:</span>
              <span className={styles.value}>{order.created_at}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>货币代码:</span>
              <span className={styles.value}>{order.currency_code}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>总金额:</span>
              <span className={styles.value}>{(order.units + order.nanos / 1e9).toFixed(2)} {order.currency_code}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>状态:</span>
              <span className={styles.value}>{order.status === 1 ? '已完成' : '进行中'}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.label}>物流状态:</span>
              <span className={styles.value}>{order.logistics_status === 0 ? '处理中' : '已发货'}</span>
            </div>
          </div>

          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>收货地址</h2>
            <p className={styles.value}>{order.street_address}</p>
            <p className={styles.value}>{order.city}, {order.state} {order.zip_code}</p>
            <p className={styles.value}>{order.country}</p>
          </div>

          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>商品详情</h2>
            {order.order_details.map((detail) => (
              <div key={detail.order_detail_id} className={styles.orderItem}>
                <div className={styles.detailRow}>
                  <span className={styles.label}>商品名称:</span>
                  <span className={styles.value}>{detail.product_name}</span>
                </div>
                <div className={styles.detailRow}>
                  <span className={styles.label}>描述:</span>
                  <span className={styles.value}>{detail.description}</span>
                </div>
                <div className={styles.detailRow}>
                  <span className={styles.label}>单价:</span>
                  <span className={styles.value}>{(detail.units + detail.nanos / 1e9).toFixed(2)} {order.currency_code}</span>
                </div>
                <hr className={styles.separator} />
              </div>
            ))}
          </div>

          <Link href="/orders" legacyBehavior>
            <a className={styles.backLink}>返回订单列表</a>
          </Link>
        </div>
      </>
    );
  }

  return (
    <div>
      <h1>加载中...</h1>
    </div>
  );
};

export default OrderDetailPage;