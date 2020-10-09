package ru.tagmeasurements.fetch_service.dtos;

import lombok.Data;

class ThresholdQ{
    public String __type;
    public double min;
    public double max;
    public double step;
    public double sample1;
    public double sample2;
}

@Data
public class TempSensorConfig {
    private String __type;
    private String email;
    private String apnsSound;
    private boolean apnsCA;
    private int apns_pause;
    private boolean send_email;
    private boolean send_tweet;
    private boolean notify_normal;
    private int temp_unit;
    private double th_low;
    private double th_high;
    private double th_window;
    private int th_low_delay;
    private int th_high_delay;
    private int interval;
    private boolean beep_pc;
    private boolean beep_pc_vibrate;
    private boolean beep_pc_tts;
    private ThresholdQ threshold_q;
}
