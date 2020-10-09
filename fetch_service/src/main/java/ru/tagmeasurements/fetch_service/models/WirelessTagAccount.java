package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.Entity;
import javax.persistence.Id;

@Data
@Entity
@NoArgsConstructor
@AllArgsConstructor
public class WirelessTagAccount {
    @Id
    private String email;
    private String password;
}
