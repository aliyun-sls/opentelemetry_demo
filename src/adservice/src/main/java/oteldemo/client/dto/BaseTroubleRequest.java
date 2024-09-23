package oteldemo.client.dto;

import java.util.HashMap;
import java.util.Map;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2022/02/08
 */
public class BaseTroubleRequest {

    private String componentName;

    private Map<String, String> params = new HashMap<>();

    public BaseTroubleRequest(String componentName) {
        this.componentName = componentName;
    }

    public BaseTroubleRequest() {
    }

    public void addParam(String key, String val) {
        this.params.put(key, val);
    }

    public String getParam(String key) {
        return params.get(key);
    }

    public String getComponentName() {
        return componentName;
    }

    public void setComponentName(String componentName) {
        this.componentName = componentName;
    }

    public Map<String, String> getParams() {
        return params;
    }

    public void setParams(Map<String, String> params) {
        this.params = params;
    }
}
