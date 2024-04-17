const hourOffsets = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12] as const;
export type HourOffset = (typeof hourOffsets)[number];
export const HourOffsets = Object.values(hourOffsets);
