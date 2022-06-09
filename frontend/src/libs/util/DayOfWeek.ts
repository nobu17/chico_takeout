export const DAY_OF_WEEK = {
  Sunday: 0,
  Monday: 1,
  TuesDay: 2,
  WednesDay: 3,
  Thursday: 4,
  Friday: 5,
  Saturday: 6,
} as const;

export type DayOfWeek = typeof DAY_OF_WEEK[keyof typeof DAY_OF_WEEK];
export const AllDayOfWeek = Object.values(DAY_OF_WEEK);

export const toString = (dayOfWeek: DayOfWeek): string => {
  switch (dayOfWeek) {
    case DAY_OF_WEEK.Sunday:
      return "日曜";
    case DAY_OF_WEEK.Monday:
      return "月曜";
    case DAY_OF_WEEK.TuesDay:
      return "火曜";
    case DAY_OF_WEEK.WednesDay:
      return "水曜";
    case DAY_OF_WEEK.Thursday:
      return "木曜";
    case DAY_OF_WEEK.Friday:
      return "金曜";
    case DAY_OF_WEEK.Saturday:
      return "土曜";
    default:
      return "";
  }
};

export const toShortString = (dayOfWeek: DayOfWeek): string => {
    switch (dayOfWeek) {
      case DAY_OF_WEEK.Sunday:
        return "日";
      case DAY_OF_WEEK.Monday:
        return "月";
      case DAY_OF_WEEK.TuesDay:
        return "火";
      case DAY_OF_WEEK.WednesDay:
        return "水";
      case DAY_OF_WEEK.Thursday:
        return "木";
      case DAY_OF_WEEK.Friday:
        return "金";
      case DAY_OF_WEEK.Saturday:
        return "土";
      default:
        return "";
    }
  };