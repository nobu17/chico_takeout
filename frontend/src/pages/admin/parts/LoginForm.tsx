import * as React from "react";
import { Container, Stack, TextField, Button } from "@mui/material";
import { SubmitHandler, useForm } from "react-hook-form";

import { isValidEmail } from "../../../libs/util/Email";
import { RequiredErrorMessage } from "../../../libs/ErrorMessages";

type AdminLoginFormProps = {
  input: LoginInput;
  onSubmit: callbackSubmit;
};
interface callbackSubmit {
  (item: LoginInput): void;
}

export type LoginInput = {
  email: string;
  password: string;
};

export default function AdminLoginForm(props: AdminLoginFormProps) {
  const {
    register,
    setError,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({ defaultValues: props.input });
  const onSubmit: SubmitHandler<LoginInput> = (data) => {
    props.onSubmit(data);
  };
  const handleEmailValidation = (email: string) => {
    const isValid = isValidEmail(email);
    if (!isValid) {
      setError("email", { message: "emailを入力してください。" });
    }
    return isValid;
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="ID"
            fullWidth
            {...register("email", {
              required: { value: true, message: RequiredErrorMessage },
              validate: handleEmailValidation,
            })}
            error={Boolean(errors.email)}
            helperText={errors.email && errors.email.message}
          />
          <TextField
            label="Password"
            fullWidth
            type="password"
            autoComplete="current-password"
            {...register("password", {
              required: { value: true, message: RequiredErrorMessage },
            })}
            error={Boolean(errors.password)}
            helperText={errors.password && errors.password.message}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            ログイン
          </Button>
        </Stack>
      </Container>
    </>
  );
}
