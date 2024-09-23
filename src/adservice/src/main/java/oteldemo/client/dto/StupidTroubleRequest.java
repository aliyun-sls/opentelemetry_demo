package oteldemo.client.dto;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2022/02/08
 */
public class StupidTroubleRequest extends BaseTroubleRequest {

    private String slowInSeconds;

    public StupidTroubleRequest(String componentName) {
        super(componentName);
    }

    public StupidTroubleRequest() {
    }

    public String getSlowInSeconds() {
        return slowInSeconds;
    }

    public void setSlowInSeconds(String slowInSeconds) {
        this.slowInSeconds = slowInSeconds;
    }
}
