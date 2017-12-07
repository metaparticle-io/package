FROM ruby:2.4-alpine

RUN mkdir -p /sinatra-app
WORKDIR /sinatra-app

COPY Gemfile /sinatra-app/
COPY Gemfile.lock /sinatra-app/
RUN bundle install

COPY . /sinatra-app/
CMD ruby app.rb
