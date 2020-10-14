package ru.tagmeasurements.fetch_service.dtos;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class MultiTagStatsRawResponse {
    private String __type;
    private List<MultiTagStats> stats;
    private List<Double> ids;
    private List<String> names;
    private List<Object> discons;
}
