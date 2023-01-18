import {
  DeepMap,
  FieldError,
  FieldValues,
  useController,
  UseControllerProps,
} from "react-hook-form";

import { Selects, SelectsProps } from "../Selects";

type OptionProps = {
  allowNoSelect?: boolean;
};

export type RhfSelectsProps<T extends FieldValues> = SelectsProps &
  UseControllerProps<T> & OptionProps;

/**
 * react-hook-formラッパー
 */
export const RhfSelects = <T extends FieldValues>(
  props: RhfSelectsProps<T>
) => {
  const { name, label, multiple, itemList, control } = props;
  const {
    field: { ref, ...rest },
    formState: { errors },
  } = useController<T>({
    name,
    control,
    rules: props.allowNoSelect ? {} : { required: { value: true, message: "項目を選択してください。" } },
  });

  return (
    <Selects
      inputRef={ref}
      {...rest}
      label={label}
      multiple={multiple}
      itemList={itemList}
      error={
        errors[name] &&
        `${(errors[name] as DeepMap<FieldValues, FieldError>).message}`
      }
    />
  );
};
