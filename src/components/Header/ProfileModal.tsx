import { Avatar, Form, Input, Message, Modal } from '@arco-design/web-react';
import React, { useEffect, useState } from 'react';

import UploadButton from '@/components/UploadButton';
import type { ProfileInput, User } from '@/utils/api';
import { updateProfile } from '@/utils/api';
import { showError } from '@/utils/feedback';

interface Props {
  visible: boolean;
  user?: User;
  onClose: () => void;
  onSaved: (user: User) => void;
}

const urlRule = { type: 'url' as const, message: '请输入以 http(s):// 开头的链接' };

const ProfileModal: React.FC<Props> = ({ visible, user, onClose, onSaved }) => {
  const [form] = Form.useForm<ProfileInput>();
  const [saving, setSaving] = useState(false);
  const avatar = Form.useWatch('avatar', form) as string | undefined;

  useEffect(() => {
    if (visible && user) {
      form.setFieldsValue({ nickname: user.nickname, avatar: user.avatar, website: user.website });
    }
  }, [visible, user, form]);

  const submit = async () => {
    const values = await form.validate();
    setSaving(true);
    try {
      const saved = await updateProfile({
        nickname: values.nickname.trim(),
        avatar: values.avatar?.trim() ?? '',
        website: values.website?.trim() ?? ''
      });
      Message.success('资料已更新！');
      onSaved(saved);
      onClose();
    } catch (err) {
      showError(err);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      title='个人资料'
      visible={visible}
      onOk={submit}
      onCancel={onClose}
      confirmLoading={saving}
      unmountOnExit
    >
      <Form form={form} layout='vertical'>
        <Form.Item label='邮箱'>
          <Input value={user?.email} disabled />
        </Form.Item>
        <Form.Item
          label='昵称'
          field='nickname'
          rules={[
            { required: true, message: '请输入昵称' },
            { maxLength: 32, message: '昵称最多 32 个字符' }
          ]}
        >
          <Input placeholder='评论区显示的昵称' />
        </Form.Item>
        <Form.Item label='头像' field='avatar' rules={[urlRule]}>
          <Input placeholder='头像图片链接，可直接上传' />
        </Form.Item>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 20 }}>
          <Avatar size={48}>{avatar ? <img src={avatar} alt='' /> : user?.nickname?.[0]}</Avatar>
          <UploadButton
            size='small'
            text='上传头像'
            onUploaded={url => form.setFieldValue('avatar', url)}
          />
        </div>
        <Form.Item label='个人网站' field='website' rules={[urlRule]}>
          <Input placeholder='https://' />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default ProfileModal;
