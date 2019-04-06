suppressPackageStartupMessages(library(tidyverse))

## state history
state_hist <- read.csv('../model/state_hist.csv')
ntemps <- length(unique(state_hist$temp))
state_plot <-
  ggplot(state_hist) +
  geom_point(aes(x = x, y = y, color = as.character(spin)), size = 0.1) +
  facet_grid(temp ~ time)
ggsave("state_hist.png", state_plot, width = 25, height = 2*ntemps, units = "cm")

## magnetization history
macro_hist <- read.csv('../model/macro_hist.csv')
macro_hist <- macro_hist %>% mutate(mag = abs(mag)) %>% gather(key = "key", value = "value", c(mag, ener))

end_time <- max(macro_hist$time)

macro_plot <-
  ggplot(macro_hist) +
  geom_point(aes(x = time, y = value), size = 0.2) +
  facet_grid(temp ~ key)
ggsave("macro_hist.png", macro_plot, width = 10, height = 2*ntemps, units = "cm")

temp_mag <-
  ggplot(macro_hist %>% filter(time == end_time)) +
  geom_point(aes(x = temp, y = value)) +
  facet_grid(. ~ key)
print(temp_mag)
