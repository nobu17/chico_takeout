import * as React from "react";
import { Stack, TextField, Button } from "@mui/material";
import { UserInfo } from "../../../hooks/UseUserInfo";
import { useForm } from "react-hook-form";
import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
} from "../../../libs/ErrorMessages";

type UserInfoInputProps = {
  userInfo: UserInfo;
  onSubmit?: submitCallback;
  onBack?: callbackBack;
};
interface submitCallback {
  (user: UserInfo): void;
}
interface callbackBack {
  (): void;
}

export default function UserInfoInput(props: UserInfoInputProps) {
  const {
    handleSubmit,
    setError,
    register,
    formState: { errors },
  } = useForm<UserInfo>({
    defaultValues: props.userInfo,
  });

  const onSubmit = (data: UserInfo) => {
    if (props.onSubmit) {
      props.onSubmit(data);
    }
  };
  const handleBack = () => {
    props.onBack?.();
  };

  const handleTelValidation = (tel: string) => {
    const reg = new RegExp("^[0-9]$");
    if (!reg.test(tel)) {
      setError("tel", { message: "数値を入力してください。" });
    }
    return true;
  };
  const handleEmailValidation = (email: string) => {
    const reg = new RegExp(
      "/^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*.)+[a-zA-Z]{2,}$/"
    );
    if (!reg.test(email)) {
      setError("email", {
        message: "正しい形式のメールアドレスを入力してください。",
      });
    }
    return true;
  };

  return (
    <Stack spacing={3}>
      <TextField
        label="氏名"
        {...register("name", {
          required: { value: true, message: RequiredErrorMessage },
          maxLength: { value: 10, message: MaxLengthErrorMessage(10) },
        })}
        error={Boolean(errors.name)}
        helperText={errors.name && errors.name.message}
      />
      <TextField
        label="E-mail"
        {...register("email", {
          required: { value: true, message: RequiredErrorMessage },
          maxLength: { value: 30, message: MaxLengthErrorMessage(30) },
          validate: handleEmailValidation,
        })}
        error={Boolean(errors.email)}
        helperText={errors.email && errors.email.message}
      />
      <TextField
        label="電話番号"
        {...register("tel", {
          required: { value: true, message: RequiredErrorMessage },
          maxLength: { value: 15, message: MaxLengthErrorMessage(10) },
          validate: handleTelValidation,
        })}
        error={Boolean(errors.tel)}
        helperText={errors.tel && errors.tel.message}
      />
      <TextField
        label="要望やメッセージ等"
        multiline
        rows={5}
        {...register("memo", {
          maxLength: { value: 500, message: MaxLengthErrorMessage(500) },
        })}
        error={Boolean(errors.memo)}
        helperText={errors.memo && errors.memo.message}
      />
      <Stack direction="row" spacing={2}>
        <Button onClick={handleSubmit(onSubmit)} variant="contained">
          次へ
        </Button>
        <Button onClick={handleBack} variant="contained" color="secondary">
          戻る
        </Button>
      </Stack>
    </Stack>
  );
}
