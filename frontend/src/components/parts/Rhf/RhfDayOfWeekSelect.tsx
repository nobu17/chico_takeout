import {
  DeepMap,
  FieldError,
  FieldValues,
  useController,
  UseControllerProps,
} from "react-hook-form";

import { DayOfWeekSelect, DayOfWeekSelectProps } from "../DayOfWeekSelect";

export type RhfDayOfWeekSelectProps<T extends FieldValues> =
  DayOfWeekSelectProps & UseControllerProps<T>;

/**
 * react-hook-formラッパー
 */
export const RhfDayOfWeekSelect = <T extends FieldValues>(
  props: RhfDayOfWeekSelectProps<T>
) => {
  const { name, label, control } = props;
  const {
    field: { ref, ...rest },
    formState: { errors },
  } = useController<T>({
    name,
    control,
    rules: { required: { value: true, message: "1つ以上の曜日を選択してください。" } },
  });

  return (
    <DayOfWeekSelect
      inputRef={ref}
      {...rest}
      label={label}
      error={
        errors[name] &&
        `${(errors[name] as DeepMap<FieldValues, FieldError>).message}`
      }
    />
  );
};
