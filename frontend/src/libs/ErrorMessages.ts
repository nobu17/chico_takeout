export const RequiredErrorMessage = "入力が必要です。";

export const MaxLengthErrorMessage = (maxLength: number) => {
  return `${maxLength}文字以下で入力してください。`;
};

export const MinErrorMessage = (min: number) => {
  return `${min}以下で入力してください。`;
};

export const MaxErrorMessage = (max: number) => {
  return `${max}以下で入力してください。`;
};
