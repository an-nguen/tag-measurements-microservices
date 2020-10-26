package ru.tagmeasurements.fetch_service.configs;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;
import ru.tagmeasurements.fetch_service.utils.LocalDateAdapter;

import java.time.LocalDate;
import java.util.concurrent.Executor;

@Configuration
@EnableAsync
public class AppConfig {
  @Bean
  @Primary
  Gson gson() {
    return new GsonBuilder()
      .setPrettyPrinting()
      .registerTypeAdapter(LocalDate.class, new LocalDateAdapter())
      .create();
  }

  @Bean
  public Executor taskExecutor() {
    ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
    executor.setCorePoolSize(16);
    executor.setMaxPoolSize(16);
    executor.setQueueCapacity(1000);
    executor.setThreadNamePrefix("WstSync-");
    executor.initialize();
    return executor;
  }
  @Bean
  public ThreadPoolTaskScheduler threadPoolTaskScheduler() {
    ThreadPoolTaskScheduler threadPoolTaskScheduler = new ThreadPoolTaskScheduler();
    threadPoolTaskScheduler.setPoolSize(3);
    return threadPoolTaskScheduler;
  }
}
