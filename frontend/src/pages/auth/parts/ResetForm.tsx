import * as React from "react";
import { Container, Stack, TextField, Button, Alert } from "@mui/material";
import { SubmitHandler, useForm } from "react-hook-form";

import { isValidEmail } from "../../../libs/util/Email";
import { RequiredErrorMessage } from "../../../libs/ErrorMessages";

type ResetFormProps = {
  input: ResetInput;
  onSubmit: callbackSubmit;
};
interface callbackSubmit {
  (data: ResetInput): void;
}

export type ResetInput = {
  email: string;
};

export default function ResetForm(props: ResetFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetInput>({ defaultValues: props.input });
  const onSubmit: SubmitHandler<ResetInput> = (data) => {
    props.onSubmit(data);
  };
  const handleEmailValidation = (email: string) => {
    const isValid = isValidEmail(email);
    if (!isValid) {
      return "正しい形式のメールアドレスを入力してください。";
    }
    return;
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="email"
            fullWidth
            {...register("email", {
              required: { value: true, message: RequiredErrorMessage },
              validate: handleEmailValidation,
            })}
            error={Boolean(errors.email)}
            helperText={errors.email && errors.email.message}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            確認メールを送信
          </Button>
        </Stack>
      </Container>
    </>
  );
}
