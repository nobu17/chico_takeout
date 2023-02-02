import * as React from "react";
import {
  Container,
  Stack,
  TextField,
  FormControlLabel,
  Checkbox,
} from "@mui/material";
import { OptionItem } from "../../../../libs/OptionItem";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
  MaxErrorMessage,
  MinErrorMessage,
} from "../../../../libs/ErrorMessages";

type OptionItemFormProps = {
  editItem: OptionItem;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (itemKind: OptionItem): void;
}
interface callbackCancel {
  (): void;
}

export default function OptionItemForm(props: OptionItemFormProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<OptionItem>({ defaultValues: props.editItem });
  const onSubmit: SubmitHandler<OptionItem> = (data) => {
    props.onSubmit(data);
  };
  const onCancel = () => {
    props.onCancel();
  };
  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 25, message: MaxLengthErrorMessage(15) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <TextField
            label="説明"
            multiline
            rows={5}
            {...register("description", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 50, message: MaxLengthErrorMessage(50) },
            })}
            error={Boolean(errors.description)}
            helperText={errors.description && errors.description.message}
          />
          <TextField
            type="number"
            label="表示順序"
            {...register("priority", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 30, message: MaxErrorMessage(30) },
            })}
            error={Boolean(errors.priority)}
            helperText={errors.priority && errors.priority.message}
          />
          <TextField
            type="number"
            label="価格"
            {...register("price", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 0, message: MinErrorMessage(0) },
              max: { value: 9999, message: MaxErrorMessage(999) },
            })}
            error={Boolean(errors.price)}
            helperText={errors.price && errors.price.message}
          />
          <FormControlLabel
            control={
              <Controller
                defaultValue={true}
                name="enabled"
                control={control}
                render={({ field: { onChange, value, ref } }) => (
                  <Checkbox
                    inputRef={ref}
                    checked={value}
                    onChange={onChange}
                  />
                )}
              />
            }
            label="有効"
          />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}
