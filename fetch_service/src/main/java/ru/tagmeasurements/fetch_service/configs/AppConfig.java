package ru.tagmeasurements.fetch_service.configs;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import ru.tagmeasurements.fetch_service.utils.LocalDateAdapter;

import java.time.LocalDate;

@Configuration
public class AppConfig {
    @Bean
    @Primary
    Gson gson() {
        return new GsonBuilder()
                .setPrettyPrinting()
                .registerTypeAdapter(LocalDate.class, new LocalDateAdapter())
                .create();
    }
}
