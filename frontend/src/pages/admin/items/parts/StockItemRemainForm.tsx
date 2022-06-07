import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { SubmitHandler, useForm, FieldError } from "react-hook-form";

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

const errorMessage = ({
  error,
  maxLength,
  min,
  max,
}: {
  name: string;
  error: FieldError | undefined;
  maxLength?: string;
  min?: string;
  max?: string;
}): string => {
  if (error?.type === "required") {
    return "入力が必要です。";
  }
  if (error?.type === "maxLength") {
    return maxLength + "文字以下にしてください。";
  }
  if (error?.type === "min") {
    return min + "以上にしてください。";
  }
  if (error?.type === "max") {
    return max + "以下にしてください。";
  }
  return "";
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
              required: true,
              min: 1,
              max: 999,
            })}
            error={Boolean(errors.remain)}
            helperText={errorMessage({
              name: "在庫数",
              error: errors.remain,
              min: "1",
              max: "999",
            })}
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
