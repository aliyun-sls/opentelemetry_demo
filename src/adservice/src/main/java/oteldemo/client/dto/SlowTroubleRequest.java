package oteldemo.client.dto;


/**
 * @author changan.zca@alibaba-inc.com
 * @date 2022/02/08
 */
public class SlowTroubleRequest extends BaseTroubleRequest {

    private String slowInSeconds;

    public SlowTroubleRequest(String componentName) {
        super(componentName);
    }

    public SlowTroubleRequest() {
    }

    public String getSlowInSeconds() {
        return slowInSeconds;
    }

    public void setSlowInSeconds(String slowInSeconds) {
        this.slowInSeconds = slowInSeconds;
    }
}
