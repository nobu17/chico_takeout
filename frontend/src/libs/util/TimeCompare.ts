export const startIsLessThanEnd = (start: string, end: string): boolean => {
  // splits number (HH:mm => xx:yy)
  const starts = parseStrTimeToNumber(start);
  const ends = parseStrTimeToNumber(end);

  const startDate = new Date(2000, 12, 12, starts[0], starts[1]);
  const endDate = new Date(2000, 12, 12, ends[0], ends[1]);

  return startDate.getTime() < endDate.getTime();
};

const parseStrTimeToNumber = (timeStr: string): number[] => {
  const splits = timeStr.split(":");
  return splits.map((m) => Number(m));
};


export const startDateIsLessThanEndDate = (start: string, end: string): boolean => {
  // splits number (yyyy/MM/dd)
  const starts = parseStrDateToNumber(start);
  const ends = parseStrDateToNumber(end);

  const startDate = new Date(starts[0], starts[1], starts[2]);
  const endDate = new Date(ends[0], ends[1], ends[2]);

  return startDate.getTime() < endDate.getTime();
};
const parseStrDateToNumber = (dateStr: string): number[] => {
  const splits = dateStr.split("/");
  return splits.map((m) => Number(m));
};


