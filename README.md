# S&P 500 Stock Tracker

이 서비스는 S&P 500 종목을 추적하고 변경사항을 구독자들에게 알려주는 REST API 서버입니다.

## 기능

1. 매일 자정에 ChatGPT API를 통해 S&P 500 종목 정보 업데이트
2. 종목 변경사항 발생 시 구독자들에게 SNS를 통해 알림
3. S&P 500 종목 목록 조회 API
4. 정량적 기준을 만족하는 종목 필터링 API

## 요구사항

- Go 1.21 이상
- AWS 계정 및 적절한 권한
- OpenAI API 키

## AWS 설정

1. S3 버킷 생성:

   - 버킷 이름: `snp500-stocks` (또는 원하는 이름)
   - 리전: `ap-northeast-2` (또는 원하는 리전)

2. SNS 토픽 생성:
   - 토픽 이름: `snp500-updates`
   - 리전: S3 버킷과 동일한 리전

## 설치 및 실행

1. 의존성 설치:

```bash
go mod download
```

2. 환경 변수 설정:

```bash
export OPENAI_API_KEY=your_api_key_here
export AWS_REGION=ap-northeast-2
export AWS_S3_BUCKET=your-bucket-name
export AWS_SNS_TOPIC_ARN=your-topic-arn
```

3. 서버 실행:

```bash
go run cmd/server/main.go
```

## API 엔드포인트

### S&P 500 종목 조회

```
GET /api/sp500
```

### 정량적 기준 종목 조회

```
GET /api/qualitative
```

### 구독 신청

```
POST /api/subscribe
Content-Type: application/json

{
    "email": "user@example.com"
}
```

## 환경 변수

- `OPENAI_API_KEY`: OpenAI API 키
- `AWS_REGION`: AWS 리전 (기본값: ap-northeast-2)
- `AWS_S3_BUCKET`: S3 버킷 이름
- `AWS_SNS_TOPIC_ARN`: SNS 토픽 ARN
- `SERVER_PORT`: 서버 포트 (기본값: 8080)
