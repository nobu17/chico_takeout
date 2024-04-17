import {
    DeepMap,
    FieldError,
    FieldValues,
    useController,
    UseControllerProps,
  } from "react-hook-form";
  
  import { HourOffsetSelect, HourOffsetSelectProps } from "../HourOffsetSelect";
  
  export type RhfHourOffsetSelectProps<T extends FieldValues> =
  HourOffsetSelectProps & UseControllerProps<T>;
  
  /**
   * react-hook-formラッパー
   */
  export const RhfHourOffsetSelect = <T extends FieldValues>(
    props: RhfHourOffsetSelectProps<T>
  ) => {
    const { name, label, control } = props;
    const {
      field: { ref, ...rest },
      formState: { errors },
    } = useController<T>({
      name,
      control,
      rules: { required: { value: true, message: "1つ以上の時間を選択してください。" } },
    });
  
    return (
      <HourOffsetSelect
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
  