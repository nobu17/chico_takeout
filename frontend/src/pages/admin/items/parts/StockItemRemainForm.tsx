import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { SubmitHandler, useForm } from "react-hook-form";
import {
  RequiredErrorMessage,
  MaxErrorMessage,
  MinErrorMessage,
} from "../../../../libs/ErrorMessages";

type StockItemRemainFormProps = {
  editItem: StockItemRemain;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (item: StockItemRemain): void;
}
interface callbackCancel {
  (): void;
}

type StockItemRemain = {
  id: string;
  name: string;
  remain: number;
};

export default function StockItemRemainForm(props: StockItemRemainFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<StockItemRemain>({ defaultValues: props.editItem });

  const onSubmit: SubmitHandler<StockItemRemain> = (data) => {
    props.onSubmit(data);
  };

  const onCancel = () => {
    props.onCancel();
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 2 }}>
        <Stack spacing={3}>
          <TextField label="名称" disabled {...register("name")} />
          <TextField
            type="number"
            label="在庫数"
            {...register("remain", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 999, message: MaxErrorMessage(10) },
            })}
            error={Boolean(errors.remain)}
            helperText={errors.remain && errors.remain.message}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            確定
          </Button>
          <Button
            color="error"
            variant="contained"
            size="large"
            onClick={onCancel}
          >
            キャンセル
          </Button>
        </Stack>
      </Container>
    </>
  );
}
