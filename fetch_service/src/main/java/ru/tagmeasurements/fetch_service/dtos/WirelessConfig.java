package ru.tagmeasurements.fetch_service.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class WirelessConfig {
    private int dataRate;
    private int activeInterval;
    @JsonProperty("Freq")
    private int freq;
    private boolean useCRC16;
    private boolean useCRC32;
    private int psid;
}
