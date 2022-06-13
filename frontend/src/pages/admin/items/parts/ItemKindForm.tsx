import * as React from "react";
import { Container, Stack, TextField } from "@mui/material";
import { ItemKind } from "../../../../libs/ItemKind";
import { SubmitHandler, useForm } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
  MaxErrorMessage,
  MinErrorMessage,
} from "../../../../libs/ErrorMessages";


type ItemKindFormProps = {
  editItem: ItemKind;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (itemKind: ItemKind): void;
}
interface callbackCancel {
  (): void;
}

export default function ItemKindForm(props: ItemKindFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ItemKind>({ defaultValues: props.editItem });
  const onSubmit: SubmitHandler<ItemKind> = (data) => {
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
              maxLength: { value: 15, message: MaxLengthErrorMessage(15) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <TextField
            type="number"
            label="表示順序"
            {...register("priority", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 10, message: MaxErrorMessage(10) },
            })}
            error={Boolean(errors.priority)}
            helperText={errors.priority && errors.priority.message}
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
