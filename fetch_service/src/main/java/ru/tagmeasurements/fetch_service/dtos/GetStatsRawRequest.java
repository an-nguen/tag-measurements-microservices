package ru.tagmeasurements.fetch_service.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDate;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class GetStatsRawRequest {
    private Long id;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate fromDate;
    @JsonFormat(pattern = "dd/MM/yyyy")
    private LocalDate toDate;
}
