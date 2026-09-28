import { Form, Input, Message, Modal } from '@arco-design/web-react';
import React, { useState } from 'react';

import { changePassword, toApiError } from '@/utils/api';
import { showError } from '@/utils/feedback';

interface Props {
  visible: boolean;
  onClose: () => void;
}

interface Values {
  current: string;
  next: string;
  confirm: string;
}

const utf8Length = (s: string) => new TextEncoder().encode(s).length;

const PasswordModal: React.FC<Props> = ({ visible, onClose }) => {
  const [form] = Form.useForm<Values>();
  const [saving, setSaving] = useState(false);

  const submit = async () => {
    const values = await form.validate();
    setSaving(true);
    try {
      await changePassword(values.current, values.next);
      Message.success('密码已修改，其它设备上的登录已失效');
      onClose();
    } catch (err) {
      if (toApiError(err).code === 'INVALID_CREDENTIALS') {
        form.setFields({ current: { error: { message: '当前密码不正确' } } });
      } else {
        showError(err);
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      title='修改密码'
      visible={visible}
      onOk={submit}
      onCancel={onClose}
      confirmLoading={saving}
      unmountOnExit
    >
      <Form form={form} layout='vertical' autoComplete='off'>
        <Form.Item label='当前密码' field='current' rules={[{ required: true, message: '请输入当前密码' }]}>
          <Input.Password autoComplete='current-password' />
        </Form.Item>
        <Form.Item
          label='新密码'
          field='next'
          rules={[
            { required: true, message: '请输入新密码' },
            {
              validator: (v: string | undefined, cb) => {
                const n = utf8Length(v ?? '');
                if (n < 8 || n > 72) cb('新密码长度需为 8-72 字节');
                else if (v === form.getFieldValue('current')) cb('新密码不能与当前密码相同');
                else cb();
              }
            }
          ]}
        >
          <Input.Password autoComplete='new-password' />
        </Form.Item>
        <Form.Item
          label='确认新密码'
          field='confirm'
          dependencies={['next']}
          rules={[
            { required: true, message: '请再次输入新密码' },
            {
              validator: (v: string | undefined, cb) =>
                v === form.getFieldValue('next') ? cb() : cb('两次输入的密码不一致')
            }
          ]}
        >
          <Input.Password autoComplete='new-password' />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default PasswordModal;
