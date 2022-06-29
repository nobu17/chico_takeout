import {
    DeepMap,
    FieldError,
    FieldValues,
    useController,
    UseControllerProps,
  } from "react-hook-form";
  
  import { Selects, SelectsProps } from "../Selects";
  
  export type RhfSelectsProps<T extends FieldValues> =
  SelectsProps & UseControllerProps<T>;
  
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
      rules: { required: { value: true, message: "項目を選択してください。" } },
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
  