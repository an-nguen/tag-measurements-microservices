package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Entity
@NoArgsConstructor
@AllArgsConstructor
public class Measurement {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private LocalDateTime date;
    private Double temperature;
    private Double humidity;
    private Double voltage;
    private Double signal;
    @Column(name = "tag_uuid")
    private UUID tagUUID;
}
