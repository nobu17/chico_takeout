export const ToDateString = (date: Date): string => {
  var y = date.getFullYear();
  var m = ("00" + (date.getMonth() + 1)).slice(-2);
  var d = ("00" + date.getDate()).slice(-2);
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
