export const toDateString = (date: Date): string => {
  const y = date.getFullYear();
  const m = ("00" + (date.getMonth() + 1)).slice(-2);
  const d = ("00" + date.getDate()).slice(-2);
  return y + "/" + m + "/" + d;
};

export const toMonthString = (date: Date): string => {
  const y = date.getFullYear();
  const m = ("00" + (date.getMonth() + 1)).slice(-2);
  return y + "/" + m;
};

export const toDate = (strDate: string): Date => {
  const splitted = strDate.split("/");
  const year = Number(splitted[0]);
  const month = Number(splitted[1]);
  const day = Number(splitted[2]);
  return new Date(year, month - 1, day);
};

export const toDateTime = (strDateTime: string): Date => {
  // only support yyyy/MM/dd HH:mm
  const [date, time] = strDateTime.split(" ");
  const splitted = date.split("/");
  const year = Number(splitted[0]);
  const month = Number(splitted[1]);
  const day = Number(splitted[2]);
  const [hour, minutes] = time.split(":");
  return new Date(year, month - 1, day, Number(hour), Number(minutes));
};

const DayOfWeeks = ["日", "月", "火", "水", "木", "金", "土"];

export const toDateTimeStrWithDayOfWeek = (date: Date): string => {
  const dateStr = toDateString(date);
  const dayOfWeek = DayOfWeeks[date.getDay()];
  const hourMin = toHourMinutesStr(date);

  return `${dateStr}(${dayOfWeek}) ${hourMin}`;
};

export const getNowDate = (addDays: number): Date => {
  const dt = getJSTNow();
  dt.setDate(dt.getDate() + addDays);
  return dt;
};

export const getNowDateByMonthOffset = (addMonth: number): Date => {
  const dt = getJSTNow();
  dt.setMonth(dt.getMonth() + addMonth);
  return dt;
};

export const isBeforeFromNow = (date: Date, addMinutesOfNow: number): boolean => {
  const now = addMinutes(getJSTNow(), addMinutesOfNow);
  return now > date;
};

export const isBeforeFromNowStr = (datetimeStr: string, addMinutesOfNow: number): boolean => {
  const datetime = toDateTime(datetimeStr);
  return isBeforeFromNow(datetime, addMinutesOfNow)
};


export const GetDateTimeFromStr = (date: string, time: string) => {
  const tempDate = toDate(date);
  const tempTime = getDateFromTimeStr(time);
  return new Date(
    tempDate.getFullYear(),
    tempDate.getMonth(),
    tempDate.getDate(),
    tempTime.getHours(),
    tempTime.getMinutes()
  );
};

export const GetTimeList = (
  startTime: string,
  endTime: string,
  offsetMinutes: number
): string[] => {
  const startDateTime = getDateFromTimeStr(startTime);
  const endDateTime = getDateFromTimeStr(endTime);

  let current = startDateTime;
  const lists: string[] = [];
  while (current < endDateTime) {
    lists.push(toHourMinutesStr(current));
    current = addMinutes(current, offsetMinutes);
  }

  return lists;
};

export const IsTimeInRange = (
  startTime: string,
  endTime: string,
  targetTime: string
): boolean => {
  const start = getDateFromTimeStr(startTime);
  const end = getDateFromTimeStr(endTime);
  const target = getDateFromTimeStr(targetTime);

  if (start <= target && end >= target) {
    return true;
  }
  return false;
};

export const convertDateTimeStrToIncludeDayOfWeeKStr = (
  dateTimeStr: string
) => {
  const date = toDateTime(dateTimeStr);

  return toDateTimeStrWithDayOfWeek(date);
};

const toHourMinutesStr = (time: Date): string => {
  const h = ("00" + time.getHours()).slice(-2);
  const m = ("00" + time.getMinutes()).slice(-2);
  return h + ":" + m;
};

const getDateFromTimeStr = (time: string): Date => {
  const [hour, minutes] = time.split(":");
  return new Date(2020, 1, 1, Number(hour), Number(minutes));
};

const addMinutes = (date: Date, minutes: number): Date => {
  return new Date(date.getTime() + minutes * 60000);
};

const getJSTNow = (): Date => {
  return new Date(Date.now() + ((new Date().getTimezoneOffset() + (9 * 60)) * 60 * 1000));
}