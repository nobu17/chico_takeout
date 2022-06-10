export const ToDateString = (date: Date): string => {
  var y = date.getFullYear();
  var m = ("00" + (date.getMonth() + 1)).slice(-2);
  var d = ("00" + date.getDate()).slice(-2);
  return y + "/" + m + "/" + d;
};

export const ToDate = (strDate: string): Date => {
  const strs = strDate.split("/");
  const year = Number(strs[0]);
  const month = Number(strs[1]);
  const day = Number(strs[2]);
  return new Date(year, month - 1, day);
};
