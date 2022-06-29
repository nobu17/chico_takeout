export const ToDateString = (date: Date): string => {
  const y = date.getFullYear();
  const m = ("00" + (date.getMonth() + 1)).slice(-2);
  const d = ("00" + date.getDate()).slice(-2);
  return y + "/" + m + "/" + d;
};

export const ToDate = (strDate: string): Date => {
  const splitted = strDate.split("/");
  const year = Number(splitted[0]);
  const month = Number(splitted[1]);
  const day = Number(splitted[2]);
  return new Date(year, month - 1, day);
};

export const GetNowDate = (addDays: number): Date => {
  const dt = new Date();
  dt.setDate(dt.getDate() + addDays);
  return dt;
};

export const IsFutureFromNow = (date: Date, offsetMinutes: number): boolean => {
  const now = addMinutes(new Date(), offsetMinutes);
  return now < date;
};

export const GetDateTimeFromStr = (date: string, time: string) => {
  const tempDate = ToDate(date);
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
