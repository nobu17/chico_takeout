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
