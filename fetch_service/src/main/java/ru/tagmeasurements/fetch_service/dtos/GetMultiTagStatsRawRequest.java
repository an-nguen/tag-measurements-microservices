package ru.tagmeasurements.fetch_service.dtos;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDate;
import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class GetMultiTagStatsRawRequest {
    public List<Long> ids;
    public LocalDate fromDate;
    public LocalDate toDate;
    public String type;
}
