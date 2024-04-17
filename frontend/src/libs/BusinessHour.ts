import { DayOfWeek } from "./util/DayOfWeek";

export type Schedules = {
  schedules: BusinessHour[];
};

export type BusinessHour = {
  id: string;
  name: string;
  start: string;
  end: string;
  weekdays: DayOfWeek[];
  enabled: boolean;
  offsetHour: number;
};
