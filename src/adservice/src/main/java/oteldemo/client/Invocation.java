package oteldemo.client;


import java.util.List;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2021/06/10
 */
public class Invocation {

    //服务
    private String service;
    //组件
    private String component;
    //方法
    private String method;

    private List<Invocation> children;

    public Invocation() {
    }

    public Invocation(String service, String component, String method, List<Invocation> children) {
        this.service = service;
        this.component = component;
        this.method = method;
        this.children = children;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getComponent() {
        return component;
    }

    public void setComponent(String component) {
        this.component = component;
    }

    public String getMethod() {
        return method;
    }

    public void setMethod(String method) {
        this.method = method;
    }

    public List<Invocation> getChildren() {
        return children;
    }

    public void setChildren(List<Invocation> children) {
        this.children = children;
    }
}
