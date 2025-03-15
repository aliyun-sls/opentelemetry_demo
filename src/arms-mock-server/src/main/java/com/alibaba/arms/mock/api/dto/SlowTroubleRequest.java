package com.alibaba.arms.mock.api.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2022/02/08
 */
@Data
@ApiModel("慢故障请求对象")
@NoArgsConstructor
public class SlowTroubleRequest extends BaseTroubleRequest {

    @ApiModelProperty("慢故障时长(秒)")
    private String slowInSeconds;

    public SlowTroubleRequest(String componentName) {
        super(componentName);
    }
}
