package oteldemo.client.dto;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2022/02/08
 */
public class ErrorTroubleRequest extends BaseTroubleRequest{
    public ErrorTroubleRequest(String componentName) {
        super(componentName);
    }

    public ErrorTroubleRequest() {
    }
}
