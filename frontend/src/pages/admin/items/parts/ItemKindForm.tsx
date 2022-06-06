import * as React from "react";
import { Button, Container, Stack, TextField } from "@mui/material";
import { ItemKind } from "../../../../libs/ItemKind";
import { SubmitHandler, useForm, FieldError } from "react-hook-form";

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

const errorMessage = ({
  name,
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
            {...register("name", { required: true, maxLength: 15 })}
            error={Boolean(errors.name)}
            helperText={errorMessage({
              name: "名称",
              error: errors.name,
              maxLength: "15",
            })}
          />
          <TextField
            type="number"
            label="表示順序"
            {...register("priority", {
              valueAsNumber: true,
              required: true,
              min: 1,
              max: 10,
            })}
            error={Boolean(errors.priority)}
            helperText={errorMessage({
              name: "表示順序",
              error: errors.priority,
              min: "1",
              max: "15",
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
