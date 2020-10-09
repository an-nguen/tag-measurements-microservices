package ru.tagmeasurements.fetch_service.dtos;

import lombok.Data;

@Data
public class TagManagerResponseData {
    private String __type;
    private Object users;
    private String name;
    private String mac;
    private String linkedToMac;
    private String notifyOfflineEmail;
    private boolean allowMore;
    private boolean selected;
    private boolean notifyOffline;
    private boolean notifyOfflinePush;
    private boolean online;
    private WirelessConfig wirelessConfig;
    private String radioId;
    private int rev;
    private int dbid;
    private String wsRoot;
    private String mStaticMAC;
}
