import {
    DeepMap,
    FieldError,
    FieldValues,
    useController,
    UseControllerProps,
  } from "react-hook-form";
  
  import { GroupSelects, GroupSelectsProps } from "../GroupSelects";
  
  type OptionProps = {
    allowNoSelect?: boolean;
  };
  
  export type RhfGroupSelectsProps<T extends FieldValues> = GroupSelectsProps &
    UseControllerProps<T> & OptionProps;
  
  /**
   * react-hook-formラッパー
   */
  export const RhfGroupSelects = <T extends FieldValues>(
    props: RhfGroupSelectsProps<T>
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
      <GroupSelects
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
  