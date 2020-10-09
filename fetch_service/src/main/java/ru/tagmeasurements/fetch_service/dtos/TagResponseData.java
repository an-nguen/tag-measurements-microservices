package ru.tagmeasurements.fetch_service.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class TagResponseData {
    private String __type;
    private int dbid;
    private String notificationJS;
    private String name;
    private String uuid;
    private String comment;
    private long slaveId;
    private int tagType;
    private Object discon;
    private Object lastComm;
    private boolean alive;
    private int signaldBm;
    private double batteryVolt;
    private boolean beeping;
    private boolean lit;
    private boolean migrationPending;
    private int beepDurationDefault;
    private int eventState;
    private int tempEventState;
    @JsonProperty("OutOfRange")
    private boolean outOfRange;
    private int tempSpurTh;
    private int lux;
    private double temperature;
    private double tempCalOffset;
    private double capCalOffset;
    private Object image_md5;
    private double cap;
    private int capRaw;
    private int az2;
    private int capEventState;
    private int lightEventState;
    private boolean shorted;
    private Object zmod;
    private Object thermostat;
    private Object playback;
    private int postBackInterval;
    private int rev;
    private int version1;
    private int freqOffset;
    private int freqCalApplied;
    private int reviveEvery;
    private int oorGrace;
    private Object tempBL;
    private Object capBL;
    private Object luxBL;
    @JsonProperty("LBTh")
    private double lBTh;
    private boolean enLBN;
    private int txpwr;
    private boolean rssiMode;
    private boolean ds18;
    private int v2flag;
    private double batteryRemaining;
}
