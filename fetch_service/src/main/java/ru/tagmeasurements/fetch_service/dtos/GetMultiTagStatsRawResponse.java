package ru.tagmeasurements.fetch_service.dtos;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class GetMultiTagStatsRawResponse {
    private MultiTagStatsRawResponse d;
}
