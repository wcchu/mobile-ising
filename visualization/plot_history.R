suppressPackageStartupMessages(library(tidyverse))

mode <- '1'

## state history
state_hist <-
  read.csv(sprintf('../model/state_hist_%s.csv', mode)) %>%
  filter(temp %% 1.0 == 0.0) %>%
  mutate(spin = as.character(spin))
ntemps <- length(unique(state_hist$temp))
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = spin), size = 0.1) +
  facet_grid(temp ~ time)
ggsave(sprintf("state_hist_%s.png", mode),
       state_plot, width = 22, height = 2*ntemps, units = "cm")

## magnetization history
macro_hist <- read.csv(sprintf('../model/macro_hist_%s.csv', mode))
macro_hist <- macro_hist %>% gather(key = "key", value = "value", c(mag, ener))
ntemps <- length(unique(macro_hist$temp))
end_time <- max(macro_hist$time)

macro_plot <-
  ggplot(macro_hist) +
  geom_point(aes(x = time, y = value), size = 0.2) +
  facet_grid(temp ~ key)
ggsave(sprintf("macro_hist_%s.png", mode),
       macro_plot, width = 10, height = 2*ntemps, units = "cm")

temp_mag <-
  ggplot(macro_hist %>% filter(time == end_time, key == "mag")) +
  geom_point(aes(x = temp, y = value)) +
  ylim(0, 1) +
  labs(x = 'Temperature', y = 'Magnetization')
ggsave(sprintf("temp_mag_%s.png", mode),
       temp_mag, width = 12, height = 10, units = "cm")

temp_ener <-
  ggplot(macro_hist %>%
           filter(key == "ener") %>%
           group_by(temp) %>%
           summarise(energy = sum(value))) +
  geom_point(aes(x = temp, y = energy)) +
  ylim(-2, 0) +
  labs(x = 'Temperature', y = 'Energy change')
ggsave(sprintf("temp_ener_%s.png", mode),
       temp_ener, width = 12, height = 10, units = "cm")
