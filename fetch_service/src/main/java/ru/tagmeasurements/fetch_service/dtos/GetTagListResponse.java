package ru.tagmeasurements.fetch_service.dtos;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class GetTagListResponse {
    private List<TagResponseData> d;
}
