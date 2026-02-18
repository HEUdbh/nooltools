<template>
  <div class="page-container">
    <div class="page-header">
      <h1>联系作者</h1>
      <p>如有问题或建议，欢迎联系</p>
    </div>
    <div class="page-content">
      <div class="contact-card">
        <div class="contact-item">
          <div class="contact-icon">📧</div>
          <div class="contact-info">
            <h3>邮箱</h3>
            <p>heuxry@outlook.com</p>
          </div>
        </div>
        <div class="contact-item">
          <div class="contact-icon">💬</div>
          <div class="contact-info">
            <h3>反馈建议</h3>
            <p>欢迎提出宝贵意见和建议</p>
          </div>
        </div>
        <div class="contact-item">
          <div class="contact-icon">⭐</div>
          <div class="contact-info">
            <h3>支持我们</h3>
            <p>感谢您的使用和支持,关注微信公众号"夜天炫安全"获取最新动态和资讯</p>
            <div class="image-container">
              <img src="/wechat.png" alt="微信公众号" style="width: 360px; margin-top: 12px;">
              <img src="/get.jpg" alt="赞助码" style="width: 240px; margin-top: 12px;">
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 联系作者页面组件
</script>

<style scoped>
.page-container {
  min-height: 100%;
  display: flex;
  flex-direction: column;
  padding: 32px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.page-header p {
  font-size: 14px;
  color: #8c92a0;
}

.page-content {
  flex: 1;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.contact-card {
  max-width: 600px;
  margin: 0 auto;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 24px;
  margin-bottom: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.contact-item:hover {
  background-color: #f0f7ff;
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.contact-icon {
  font-size: 32px;
  margin-right: 20px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ffffff;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.contact-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 4px;
}

.contact-info p {
  font-size: 14px;
  color: #5a6c7d;
}

.image-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  margin-top: 12px;
}

.image-container img {
  flex-shrink: 0;
}
</style>