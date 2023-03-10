import * as React from "react";
import { Container, Stack, TextField, Button } from "@mui/material";
import { Message } from "../../../../libs/Messages";
import { SubmitHandler, useForm } from "react-hook-form";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
} from "../../../../libs/ErrorMessages";

type MessageFormProps = {
  editItem?: Message;
  onSubmit: callbackSubmit;
};
interface callbackSubmit {
  (message: Message): void;
}

const convertItem = (editItem: Message | undefined) => {
  if (!editItem) {
    return { id: "", content: "メッセージ種別を選択してください。" };
  }
  return JSON.parse(JSON.stringify(editItem));
};

const validSelectedItem = (message: Message | undefined): boolean => {
  if (!message || !message.id) {
    return false;
  }
  return true;
};

export default function MessageForm(props: MessageFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isValid },
  } = useForm<Message>({
    defaultValues: convertItem(props.editItem),
    mode: "all",
    reValidateMode: "onChange",
  });
  const onSubmit: SubmitHandler<Message> = (data) => {
    props.onSubmit(data);
  };

  React.useEffect(() => {
    const item = convertItem(props.editItem);
    setValue("id", item.id);
    setValue("content", item.content);
  }, [props.editItem]);

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <input
            type="hidden"
            {...register("id", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 3, message: MaxLengthErrorMessage(3) },
            })}
          />
          <TextField
            label="メッセージ内容(1000文字まで)"
            multiline
            disabled={!validSelectedItem(props.editItem)}
            rows={7}
            {...register("content", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 1000, message: MaxLengthErrorMessage(1000) },
            })}
            error={Boolean(errors.content)}
            helperText={errors.content && errors.content.message}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            disabled={!isValid}
            onClick={handleSubmit(onSubmit)}
          >
            確定
          </Button>
        </Stack>
      </Container>
    </>
  );
}
