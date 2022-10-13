export const toShortenString = (target: string, maxLength: number): string => {
  if (target.length <= maxLength) {
    return target;
  }
  return target.substr(0, maxLength) + "...";
};
