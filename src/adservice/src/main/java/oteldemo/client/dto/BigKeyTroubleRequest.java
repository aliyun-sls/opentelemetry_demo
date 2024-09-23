package oteldemo.client.dto;

public class BigKeyTroubleRequest extends BaseTroubleRequest{
    private String size;
    // key过多 或者 value值过大
    private String type;

    public BigKeyTroubleRequest(String componentName) {
        super(componentName);
    }

    public String getSize() {
        return size;
    }

    public void setSize(String size) {
        this.size = size;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public BigKeyTroubleRequest() {
    }
}
