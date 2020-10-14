package ru.tagmeasurements.fetch_service.dtos;

import lombok.Data;

import java.util.List;

@Data
public class MultiTagStats {
    public String date;
    public List<Long> ids;
    public List<List<Double>> values;
    public Object values_base64;
    public List<List<Long>> tods;
    public Object tods_base64;
}
