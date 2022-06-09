import { DeepMap, FieldError, FieldValues, useController, UseControllerProps } from 'react-hook-form';

import { TimeSelect, TimeSelectProps } from '../TimeSelect';

export type RhfTimeSelectProps<T extends FieldValues> = TimeSelectProps & UseControllerProps<T>;

/**
 * react-hook-formラッパー
 */
export const RhfTimeSelect = <T extends FieldValues>(props: RhfTimeSelectProps<T>) => {
  const { name, label, timeList, control, } = props;
  const {
    field: { ref, ...rest },
    formState: { errors },
  } = useController<T>({ name, control });

  return (
    <TimeSelect
      inputRef={ref}
      {...rest}
      label={label}
      timeList={timeList}
      error={errors[name] && `${(errors[name] as DeepMap<FieldValues, FieldError>).message}`}
    />
  );
};